// Generated automatically from software.amazon.awssdk.services.s3.model.ObjectLockRule for testing purposes

package software.amazon.awssdk.services.s3.model;

import java.io.Serializable;
import java.util.List;
import java.util.Optional;
import java.util.function.Consumer;
import software.amazon.awssdk.core.SdkField;
import software.amazon.awssdk.core.SdkPojo;
import software.amazon.awssdk.services.s3.model.DefaultRetention;
import software.amazon.awssdk.utils.builder.CopyableBuilder;
import software.amazon.awssdk.utils.builder.ToCopyableBuilder;

public class ObjectLockRule implements SdkPojo, Serializable, ToCopyableBuilder<ObjectLockRule.Builder, ObjectLockRule>
{
    protected ObjectLockRule() {}
    public ObjectLockRule.Builder toBuilder(){ return null; }
    public final <T> java.util.Optional<T> getValueForField(String p0, java.lang.Class<T> p1){ return null; }
    public final DefaultRetention defaultRetention(){ return null; }
    public final List<SdkField<? extends Object>> sdkFields(){ return null; }
    public final String toString(){ return null; }
    public final boolean equals(Object p0){ return false; }
    public final boolean equalsBySdkFields(Object p0){ return false; }
    public final int hashCode(){ return 0; }
    public static ObjectLockRule.Builder builder(){ return null; }
    public static java.lang.Class<? extends ObjectLockRule.Builder> serializableBuilderClass(){ return null; }
    static public interface Builder extends CopyableBuilder<ObjectLockRule.Builder, ObjectLockRule>, SdkPojo
    {
        ObjectLockRule.Builder defaultRetention(DefaultRetention p0);
        default ObjectLockRule.Builder defaultRetention(java.util.function.Consumer<DefaultRetention.Builder> p0){ return null; }
    }
}
